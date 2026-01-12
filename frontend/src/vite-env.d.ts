/// <reference types="vite/client" />

declare module '*.vue' {
    import type { DefineComponent } from 'vue'
    const component: DefineComponent<{}, {}, any>
    export default component
}

// Wails runtime types
interface WailsRuntime {
    WindowMinimise(): void
    WindowMaximise(): void
    WindowUnmaximise(): void
    WindowToggleMaximise(): void
    WindowIsMaximised(): Promise<boolean>
    WindowIsMinimised(): Promise<boolean>
    WindowClose(): void
    BrowserOpenURL(url: string): void
    Environment(): Promise<{
        buildType: string
        platform: string
        arch: string
    }>
    Quit(): void
    ClipboardSetText(text: string): void
    ClipboardGetText(): Promise<string>
    LogError(message: string): void
    LogInfo(message: string): void
    LogDebug(message: string): void
    EventsOn(event: string, callback: (...args: any[]) => void): void
    EventsOff(event: string): void
    EventsOnMultiple(event: string, callback: (...args: any[]) => void, count: number): void
    EventsOnce(event: string, callback: (...args: any[]) => void): void
}

// Wails Go bindings types
interface WailsGoBindings {
    main: {
        App: {
            LaunchGame(): Promise<void>
            [key: string]: (...args: any[]) => Promise<any>
        }
    }
}

declare global {
    interface Window {
        runtime: WailsRuntime
        go: WailsGoBindings
    }
}

export {}
