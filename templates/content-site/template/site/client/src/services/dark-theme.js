import { storage as appStorage } from '../storages/app';

const COOKIE_KEY = 't';
const COOKIE_VALUE = 'd';

const CONTAINER_CLASS = 'dark';
const DARK_THEME_COLOR = '#2b2b2b';

export default class DarkTheme {
    constructor( DI, container ) {
        DI.bind( 'dark-theme', this );

        this.container = container;

        /** @type {Cookies} */
        this.cookies = DI.resolve( 'cookies' );

        // Get current state from cookies
        const enabled = COOKIE_VALUE === this.cookies.get( COOKIE_KEY );

        // Setup to app storage
        const app = appStorage( DI );
        app.darkThemeEnabled = enabled;
        app.sync();
        this.app = app;

        if ( !IS_SSR ) {

            // Get meta tag from head
            this.themeColor = document.head.querySelector( 'meta[name="theme-color"]' );
        }

        if ( !PRODUCTION ) {

            // For dev only env force set class, because dev-server don't add class to container
            if ( enabled ) {
                this.container.classList.add( CONTAINER_CLASS );
                this.themeColor.content = DARK_THEME_COLOR;
            }
        }
    }

    toggle() {
        const enabled = !this.app.darkThemeEnabled;

        // Set to store
        this.app.darkThemeEnabled = enabled;
        this.app.sync();

        // Save in cookies
        if ( enabled ) {
            this.cookies.set( COOKIE_KEY, COOKIE_VALUE );
            this.themeColor.content = DARK_THEME_COLOR;
        } else {
            this.cookies.delete( COOKIE_KEY );
            this.themeColor.content = '';
        }

        // Toggle class
        this.container.classList.toggle( CONTAINER_CLASS );
    }
}
