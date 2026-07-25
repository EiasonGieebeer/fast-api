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
import assert from 'node:assert/strict'
import { describe, test } from 'node:test'

import { isCurrentBuildBasedOnRelease } from '../version'

describe('update checker version comparison', () => {
  test('accepts an exact tagged release as current', () => {
    assert.equal(
      isCurrentBuildBasedOnRelease('v1.0.0-rc.21', 'v1.0.0-rc.21'),
      true
    )
  })

  test('accepts custom build metadata based on the latest release', () => {
    assert.equal(
      isCurrentBuildBasedOnRelease(
        'v1.0.0-rc.21+my-custom.fcb17e8b8',
        'v1.0.0-rc.21'
      ),
      true
    )
  })

  test('rejects a custom build based on an older release', () => {
    assert.equal(
      isCurrentBuildBasedOnRelease(
        'v1.0.0-rc.16+my-custom.53f6cfce0',
        'v1.0.0-rc.21'
      ),
      false
    )
  })
})
