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
import { after, describe, test } from 'node:test'

import { Window } from 'happy-dom'

const domWindow = new Window()
const domGlobals = [
  'window',
  'document',
  'navigator',
  'HTMLElement',
  'HTMLInputElement',
  'HTMLButtonElement',
  'HTMLAnchorElement',
  'SVGElement',
  'Node',
  'Element',
  'Event',
  'CustomEvent',
  'MutationObserver',
  'ResizeObserver',
  'requestAnimationFrame',
  'cancelAnimationFrame',
  'getComputedStyle',
] as const

for (const key of domGlobals) {
  Object.defineProperty(globalThis, key, {
    configurable: true,
    value: domWindow[key],
  })
}

const { act } = await import('react')
const { createRoot } = await import('react-dom/client')
const { createInstance } = await import('i18next')
const { I18nextProvider, initReactI18next } = await import('react-i18next')
const { RechargeFormCard } = await import('../recharge-form-card')

const i18n = createInstance()
await i18n.use(initReactI18next).init({
  lng: 'en',
  resources: {
    en: {
      translation: {
        'Buy codes at Liandong Store (Recommended)':
          'Buy codes at Liandong Store (Recommended)',
        'Pay there, then redeem below': 'Pay there, then redeem below',
        'Get one here': 'Get one here',
      },
    },
  },
})

const reactTestGlobals = globalThis as typeof globalThis & {
  IS_REACT_ACT_ENVIRONMENT?: boolean
}
reactTestGlobals.IS_REACT_ACT_ENVIRONMENT = true

describe('wallet recharge form', () => {
  after(() => {
    domWindow.close()
  })

  test('shows the configured redemption store before Alipay and keeps the lower link', async () => {
    const container = document.createElement('div')
    document.body.append(container)
    const root = createRoot(container)
    const topupLink = 'https://pay.ldxp.cn/shop/fastapi'

    await act(async () => {
      root.render(
        <I18nextProvider i18n={i18n}>
          <RechargeFormCard
            topupInfo={{
              enable_online_topup: true,
              enable_stripe_topup: false,
              pay_methods: [{ name: 'Alipay', type: 'alipay' }],
              min_topup: 1,
              stripe_min_topup: 1,
              amount_options: [10],
              discount: {},
              enable_redemption: true,
            }}
            presetAmounts={[{ value: 10 }]}
            selectedPreset={10}
            onSelectPreset={() => {}}
            topupAmount={10}
            onTopupAmountChange={() => {}}
            paymentAmount={10}
            calculating={false}
            onPaymentMethodSelect={() => {}}
            paymentLoading={null}
            redemptionCode=''
            onRedemptionCodeChange={() => {}}
            onRedeem={() => {}}
            redeeming={false}
            topupLink={topupLink}
          />
        </I18nextProvider>
      )
    })

    const storeLinks = container.querySelectorAll<HTMLAnchorElement>(
      `a[href="${topupLink}"]`
    )
    assert.equal(storeLinks.length, 2)
    assert.equal(storeLinks[0]?.target, '_blank')
    assert.equal(storeLinks[0]?.rel, 'noopener noreferrer')
    assert.match(
      storeLinks[0]?.textContent ?? '',
      /Buy codes at Liandong Store/
    )
    assert.match(storeLinks[0]?.textContent ?? '', /Pay there, then redeem below/)

    const alipayButton = [...container.querySelectorAll('button')].find(
      (button) => button.textContent?.includes('Alipay')
    )
    assert.ok(alipayButton)
    assert.ok(
      (storeLinks[0]?.compareDocumentPosition(alipayButton) ?? 0) &
        Node.DOCUMENT_POSITION_FOLLOWING
    )
    assert.match(storeLinks[1]?.textContent ?? '', /Get one here/)

    await act(async () => root.unmount())
    container.remove()
  })
})
