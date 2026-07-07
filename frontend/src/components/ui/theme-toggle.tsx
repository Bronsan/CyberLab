'use client'

import { useTheme } from '@/providers/theme-provider'
import { Sun, Moon, Monitor } from 'lucide-react'

export function ThemeToggle() {
  const { theme, setTheme, resolvedTheme } = useTheme()

  const cycleTheme = () => {
    if (theme === 'light') setTheme('dark')
    else if (theme === 'dark') setTheme('system')
    else setTheme('light')
  }

  return (
    <button
      onClick={cycleTheme}
      className="relative p-2 rounded-lg hover:bg-muted transition-colors"
      title={`Current: ${theme}`}
    >
      {theme === 'light' && <Sun className="h-4 w-4 text-amber-500" />}
      {theme === 'dark' && <Moon className="h-4 w-4 text-blue-400" />}
      {theme === 'system' && <Monitor className="h-4 w-4 text-muted-foreground" />}
    </button>
  )
}
