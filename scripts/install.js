const { name, goVersion: version } = require('../package.json')
const { createHash } = require('crypto')
let { chmod, mkdir, writeFile } = require('fs')
const { concat: downloadBuffer } = require('simple-get')
const { resolve } = require('path')
const { promisify } = require('util')
const { t: openTar } = require('tar')
const { fromBuffer: openZip } = require('yauzl')
const streamToBuffer = require('gather-stream')
const stream = require('stream')
mkdir = promisify(mkdir)
chmod = promisify(chmod)
writeFile = promisify(writeFile)

const platformAliases = {
  darwin: 'darwin',
  linux: 'linux',
  win32: 'windows'
}
const archAliases = {
  arm: 'armv6',
  arm64: 'arm64',
  x64: 'x86_64'
}
const formats = {
  windows: 'zip'
}
const extensions = {
  windows: '.exe'
}
const decompressors = {
  windows: decompressZip
}

let { platform, arch } = process
platform = platformAliases[platform]
if (!platform) fail(`Unsupported platform: ${process.platform}.`)
arch = archAliases[arch]
if (!arch) fail(`Unsupported architecture: ${process.arch}.`)
const format = formats[platform] || 'tar.gz'
const extension = extensions[platform] || ''
const decompress = decompressors[platform] || decompressTar

const releaseURL = `https://github.com/prantlf/saz-tools/releases/download/v${version}/`
const checksumURL = `${releaseURL}checksums.txt`
const archiveName = `${name}_${version}_${platform}_${arch}.${format}`
const archiveURL = `${releaseURL}${archiveName}`

if (name === 'saz-tools') {
  console.log('Skipping installation of the binary within the original package.')
  process.exit(0)
}
if (process.env.DYNO !== undefined) {
  console.log('Skipping installation of the binary on Heroku.')
  process.exit(0)
}

Promise
  .all([downloadFile(archiveURL), downloadFile(checksumURL)])
  .then(installBinaries)
  .catch(({ message }) => fail(message))

async function installBinaries ([archive, checksums]) {
  await checkHash(archive, checksums)
  const [sazdump, sazserve] = await Promise.all([
    decompress(archive, 'sazdump'),
    decompress(archive, 'sazserve')])
  const directory = getBinaryDirectory()
  await Promise.all([
    saveFile('sazdump', directory, sazdump),
    saveFile('sazserve', directory, sazserve)])
  console.log()
}

function downloadFile (url) {
  return new Promise((resolve, reject) => {
    console.log(`Downloading ${url}.`)
    downloadBuffer(url, (error, response, data) => {
      if (error) return reject(error)
      if (response.statusCode !== 200) {
        return reject(new Error(
          `Downloading ${url} failed: ${response.statusCode} ${response.statusMessage}.`))
      }
      resolve(data)
    })
  })
}

function findExpectedHash (content) {
  for (const line of content.split(/\r?\n/)) {
    const [expectedHash, sourceFile] = line.split(/\s+/)
    if (sourceFile === archiveName) return expectedHash
  }
  throw new Error('Checksum for the downloaded archive not found.')
}

function computeActualHash (content) {
  return new Promise((resolve, reject) => {
    const hash = createHash('sha256')
    const input = new stream.PassThrough()
    input
      .on('readable', () => {
        const data = input.read()
        if (data) hash.update(data)
        else resolve(hash.digest('hex'))
      })
      .on('error', reject)
      .end(content)
  })
}

async function checkHash (archive, checksums) {
  console.log(`Checking content of ${archiveName}.`)
  const expectedHash = findExpectedHash(checksums.toString())
  const actualHash = await computeActualHash(archive)
  if (expectedHash !== actualHash) { throw new Error('Checksum of the downloaded archive differed.') }
}

async function decompressTar (archive, name) {
  return new Promise((resolve, reject) => {
    const data = []
    name += extension
    openTar()
      .on('entry', entry => {
        if (entry.path === name) {
          entry
            .on('error', reject)
            .on('data', chunk => data.push(chunk))
        }
      })
      .on('error', reject)
      .on('close', () => {
        if (!data.length) return reject(new Error(`File .../${name} not found.`))
        resolve(Buffer.concat(data))
      })
      .end(archive)
  })
}

async function decompressZip (archive, name) {
  return new Promise((resolve, reject) => {
    let buffer
    name += extension
    openZip(archive, { lazyEntries: true }, (error, zip) => {
      if (error) reject(error)
      else {
        zip
          .on('entry', entry => {
            if (entry.fileName === name) {
              zip.openReadStream(entry, (error, stream) => {
                if (error) return reject(error)
                stream.on('end', () => zip.readEntry())
                stream.pipe(streamToBuffer((error, result) => {
                  if (error) reject(error)
                  else buffer = result
                }))
              })
            } else {
              zip.readEntry()
            }
          })
          .on('end', () => {
            if (!buffer) reject(new Error(`File .../${name} not found.`))
            else resolve(buffer)
          })
          .readEntry()
      }
    })
  })
}

function getBinaryDirectory () {
  const root = resolve(process.execPath, '../..')
  return __dirname.startsWith(root) ? `${root}/bin` : `${__dirname}/../node_modules/.bin`
}

async function saveFile (name, directory, content) {
  const path = `${directory}/${name}${extension}`
  console.log(`Writing ${path}.`)
  await mkdir(directory, { recursive: true })
  await writeFile(path, content)
  await chmod(path, 0o755)
}

function fail (message) {
  console.error(message)
  console.log()
  process.exit(1)
}
