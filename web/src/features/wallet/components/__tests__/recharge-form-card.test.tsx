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
import { render, screen } from '@testing-library/react'
import { describe, expect, test } from 'vitest'

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

describe('wallet recharge form', () => {
  test('shows the configured redemption store before Alipay and keeps the lower link', () => {
    const topupLink = 'https://pay.ldxp.cn/shop/fastapi'

    const { container } = render(
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

    const storeLinks = container.querySelectorAll<HTMLAnchorElement>(
      `a[href="${topupLink}"]`
    )
    expect(storeLinks).toHaveLength(2)
    expect(storeLinks[0]).toHaveAttribute('target', '_blank')
    expect(storeLinks[0]).toHaveAttribute('rel', 'noopener noreferrer')
    expect(storeLinks[0]).toHaveTextContent('Buy codes at Liandong Store')
    expect(storeLinks[0]).toHaveTextContent('Pay there, then redeem below')

    const alipayButton = screen.getByRole('button', { name: /Alipay/ })
    expect(
      (storeLinks[0]?.compareDocumentPosition(alipayButton) ?? 0) &
        Node.DOCUMENT_POSITION_FOLLOWING
    ).not.toBe(0)
    expect(storeLinks[1]).toHaveTextContent('Get one here')
  })
})
