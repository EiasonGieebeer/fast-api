import { useEffect, useState } from 'react'
import { ExternalLink, Loader2 } from 'lucide-react'
import { useTranslation } from 'react-i18next'
import { toast } from 'sonner'
import { Button } from '@/components/ui/button'
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog'
import { Input } from '@/components/ui/input'
import { tryPrettyJson } from '@/lib/utils'
import { completeCodexOAuth, startCodexOAuth } from '../../api'

type Props = {
  open: boolean
  onOpenChange: (open: boolean) => void
  onKeyGenerated: (key: string) => void
}

export function CodexOAuthDialog({
  open,
  onOpenChange,
  onKeyGenerated,
}: Props) {
  const { t } = useTranslation()
  const [authorizeUrl, setAuthorizeUrl] = useState('')
  const [callbackUrl, setCallbackUrl] = useState('')
  const [loading, setLoading] = useState(false)

  useEffect(() => {
    if (!open) {
      setAuthorizeUrl('')
      setCallbackUrl('')
      setLoading(false)
    }
  }, [open])

  const handleStart = async () => {
    setLoading(true)
    try {
      const res = await startCodexOAuth()
      const url = res.data?.authorize_url
      if (!res.success || !url) throw new Error(res.message || 'OAuth failed')
      setAuthorizeUrl(url)
      window.open(url, '_blank', 'noopener,noreferrer')
    } catch (error) {
      toast.error(error instanceof Error ? error.message : t('OAuth failed'))
    } finally {
      setLoading(false)
    }
  }

  const handleComplete = async () => {
    setLoading(true)
    try {
      const res = await completeCodexOAuth(callbackUrl.trim())
      const key = res.data?.key
      if (!res.success || !key) throw new Error(res.message || 'OAuth failed')
      onKeyGenerated(tryPrettyJson(key))
      toast.success(t('Credential generated'))
      onOpenChange(false)
    } catch (error) {
      toast.error(error instanceof Error ? error.message : t('OAuth failed'))
    } finally {
      setLoading(false)
    }
  }

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className='sm:max-w-2xl'>
        <DialogHeader>
          <DialogTitle>{t('Codex Authorization')}</DialogTitle>
          <DialogDescription>
            {t(
              'Open the authorization page, then paste the full localhost callback URL below.'
            )}
          </DialogDescription>
        </DialogHeader>
        <div className='space-y-3'>
          <div className='flex gap-2'>
            <Button onClick={handleStart} disabled={loading}>
              {loading ? (
                <Loader2 className='mr-2 h-4 w-4 animate-spin' />
              ) : (
                <ExternalLink className='mr-2 h-4 w-4' />
              )}
              {t('Open authorization page')}
            </Button>
            {authorizeUrl && (
              <Button
                type='button'
                variant='outline'
                onClick={() => navigator.clipboard.writeText(authorizeUrl)}
              >
                {t('Copy authorization link')}
              </Button>
            )}
          </div>
          <Input
            value={callbackUrl}
            onChange={(event) => setCallbackUrl(event.target.value)}
            placeholder={t('Paste the full callback URL')}
            autoComplete='off'
          />
        </div>
        <DialogFooter>
          <Button
            type='button'
            variant='outline'
            onClick={() => onOpenChange(false)}
          >
            {t('Cancel')}
          </Button>
          <Button
            onClick={handleComplete}
            disabled={loading || !callbackUrl.trim()}
          >
            {t('Generate credential')}
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  )
}
