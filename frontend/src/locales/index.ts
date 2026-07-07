import { en } from './en'
import { zh } from './zh'
import { ja } from './ja'

export type Locale = 'en' | 'zh' | 'ja'
export type Translation = typeof en

export const translations: Record<Locale, Translation> = { en, zh, ja }

export const localeNames: Record<Locale, string> = {
  en: 'English',
  zh: '中文',
  ja: '日本語',
}

export const localeFlags: Record<Locale, string> = {
  en: '🇺🇸',
  zh: '🇨🇳',
  ja: '🇯🇵',
}

export const defaultLocale: Locale = 'en'
