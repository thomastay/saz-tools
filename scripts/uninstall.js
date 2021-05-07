let { unlink } = require('fs')
const { join, resolve } = require('path')
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
  const binaryDirectory = getBinaryDirectory()
  const manualDirectory = getManualDirectory()
  await Promise.all([
    deleteBinary('sazdump', binaryDirectory), deleteBinary('sazserve', binaryDirectory),
    deleteManual('sazdump', manualDirectory), deleteManual('sazserve', manualDirectory)])
  console.log()
}

function getBinaryDirectory () {
  const root = resolve(process.execPath, '../..')
  return __dirname.startsWith(root) ? join(root, 'bin') : join(__dirname, '../node_modules/.bin')
}

function getManualDirectory () {
  const root = resolve(process.execPath, '../..')
  return __dirname.startsWith(root) && `${root}/share/man/man1`
}

async function deleteBinary (name, directory) {
  const path = `${directory}/${name}${extension}`
  console.log(`Deleting ${path}.`)
  await unlink(path)
}

async function deleteManual (name, directory) {
  if (directory) {
    const path = `${directory}/${name}.1`
    console.log(`Deleting ${path}.`)
    await unlink(path)
  }
}

function fail (message) {
  console.error(message)
  console.log()
  process.exit(1)
}
