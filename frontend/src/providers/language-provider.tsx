'use client'

import { createContext, useContext, useState, useEffect, useCallback, type ReactNode } from 'react'
import { type Locale, type Translation, translations, defaultLocale } from '@/locales'

interface LanguageContextType {
  locale: Locale
  t: Translation
  setLocale: (locale: Locale) => void
}

const LanguageContext = createContext<LanguageContextType>({
  locale: defaultLocale,
  t: translations[defaultLocale],
  setLocale: () => {},
})

export function useTranslation() {
  return useContext(LanguageContext)
}

export function LanguageProvider({ children }: { children: ReactNode }) {
  const [locale, setLocaleState] = useState<Locale>(defaultLocale)

  useEffect(() => {
    // Load saved locale
    const saved = localStorage.getItem('cyberlab-locale') as Locale | null
    if (saved && translations[saved]) {
      setLocaleState(saved)
    } else {
      // Detect browser language
      const browserLang = navigator.language.split('-')[0]
      if (browserLang === 'zh') setLocaleState('zh')
      else if (browserLang === 'ja') setLocaleState('ja')
    }
  }, [])

  const setLocale = useCallback((locale: Locale) => {
    setLocaleState(locale)
    localStorage.setItem('cyberlab-locale', locale)
    document.documentElement.lang = locale
  }, [])

  const t = translations[locale]

  return (
    <LanguageContext.Provider value={{ locale, t, setLocale }}>
      {children}
    </LanguageContext.Provider>
  )
}
