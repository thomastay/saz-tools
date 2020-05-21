let { unlink } = require('fs')
const { resolve } = require('path')
const { promisify } = require('util')
unlink = promisify(unlink)

const platformAliases = {
  darwin: 'darwin',
  linux: 'linux',
  win32: 'windows'
}
const extensions = {
  windows: '.exe'
}

let { platform } = process
platform = platformAliases[platform]
if (!platform) fail(`Unsupported platform: ${process.platform}.`)
const extension = extensions[platform] || ''

uninstallBinaries()
  .catch(({ message }) => fail(message))

async function uninstallBinaries () {
  const directory = getBinaryDirectory()
  await Promise.all([
    deleteFile('sazdump', directory), deleteFile('sazserve', directory)])
  console.log()
}

function getBinaryDirectory () {
  const root = resolve(process.execPath, '../..')
  return __dirname.startsWith(root) ? `${root}/bin` : `${__dirname}/node_modules/.bin`
}

async function deleteFile (name, directory) {
  const path = `${directory}/${name}${extension}`
  console.log(`Deleting ${path}.`)
  await unlink(path)
}

function fail (message) {
  console.error(message)
  console.log()
  process.exit(1)
}
