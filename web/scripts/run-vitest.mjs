/*
Copyright (C) 2023-2026 QuantumNous

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU Affero General Public License as
published by the Free Software Foundation, either version 3 of the
License, or (at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
GNU Affero General Public License for more details.

You should have received a copy of the GNU Affero General Public License
along with this program. If not, see <https://www.gnu.org/licenses/>.

For commercial licensing, please contact support@quantumnous.com
*/
import { spawnSync } from 'node:child_process'

const flagSupport = spawnSync(
  'node',
  [
    '-p',
    "process.allowedNodeEnvironmentFlags.has('--no-experimental-webstorage')",
  ],
  { encoding: 'utf8' }
)
if (flagSupport.error) {
  throw flagSupport.error
}

const childEnv = { ...process.env }
if (flagSupport.status === 0 && flagSupport.stdout.trim() === 'true') {
  childEnv.NODE_OPTIONS = [
    childEnv.NODE_OPTIONS,
    '--no-experimental-webstorage',
  ]
    .filter(Boolean)
    .join(' ')
}

const result = spawnSync(
  'node',
  ['./node_modules/vitest/vitest.mjs', 'run', ...process.argv.slice(2)],
  { env: childEnv, stdio: 'inherit' }
)

if (result.error) {
  throw result.error
}

process.exit(result.status ?? 1)
