import { start, createDI } from 'sham-ui';
import pretty from 'pretty';
import { renderAPP, toHTML } from '../../src/ssr';
import { setup as setupRehydrator } from 'sham-ui-ssr/lib/rehydrator';
import initializer from '../../src/initializers/main';

export const app = {
    async start( DI, waitRendering = true ) {
        if ( !DI ) {
            DI = createDI();
        }
        DI
            .bind( 'location:origin', document.location.origin + '/' )
            .bind( 'location:href', document.location.href )
            .bind( 'document', document )
        ;
        initializer( DI, document.querySelector( 'body' ) );
        start( DI );
        if ( waitRendering ) {
            await this.waitRendering();
        }
    },

    async ssrAndCheck( url, cookies = '' ) {
        const data = await renderAPP(
            'http://localhost:3002/api/',
            document.location.origin + '/',
            url,
            cookies
        );
        const html = toHTML( data );
        expect(
            pretty( html, {
                inline: [ 'code', 'pre', 'em', 'strong', 'span' ]
            } ),
        ).toMatchSnapshot();
    },

    async ssrAndRehydrate( url, DI, cookies = '' ) {
        setupSSR();
        const data = await renderAPP(
            'http://localhost:3002/api/',
            document.location.origin + '/',
            url,
            cookies
        );
        setup();
        document.querySelector( 'html' ).innerHTML = toHTML( data );
        history.pushState( {}, '', url );
        if ( !DI ) {
            DI = createDI();
        }
        DI
            .bind( 'location:origin', document.location.origin + '/' )
            .bind( 'location:href', document.location.href )
            .bind( 'document', document )
        ;
        const disableRehydrating = setupRehydrator(
            DI,
            JSON.parse( data.data )
        );
        await initializer(
            DI,
            document.querySelector( 'body' )
        );
        start( DI );
        disableRehydrating();
    },

    async waitRendering() {
        await new Promise( resolve => setImmediate( resolve ) );
    },
    click( selector ) {
        document.querySelector( selector ).click();
    },
    form: {
        fill( field, value ) {
            this.fillBySelector( `[name="${field}"]`, value );
        },
        fillBySelector( selector, value ) {
            const element = document.querySelector( selector );
            element.value = value;
            element.dispatchEvent( new Event( 'input' ) );
            element.dispatchEvent( new Event( 'change' ) );
        },
        async submit() {
            app.click( '[type="submit"]' );
            await app.waitRendering();
        }
    },
    checkBody() {
        expect(
            pretty( document.querySelector( 'body' ).innerHTML, {
                inline: [ 'code', 'pre', 'em', 'strong', 'span' ]
            } ),
        ).toMatchSnapshot();
    },
    checkMainPanel() {
        expect(
            pretty( document.querySelector( '.mainpanel' ).innerHTML, {
                inline: [ 'code', 'pre', 'em', 'strong', 'span' ]
            } ),
        ).toMatchSnapshot();
    }
};

function setupRAF() {
    window.requestAnimationFrame = ( cb ) => {
        setImmediate( cb );
    };
}

function clearBody() {
    document.querySelector( 'body' ).innerHTML = '';
}

function setupRouter() {
    delete window.__NAVIGO_WINDOW_LOCATION_MOCK__;
    history.pushState( {}, '', '' );
}

export function setupSSR() {
    Object.defineProperty( window, 'IS_SSR', {
        value: true,
        configurable: true,
        writable: true
    } );
}

export default function setup( DI ) {
    if ( !DI ) {
        DI = createDI();
    }
    setupRAF();
    clearBody();
    setupRouter();
    Object.defineProperty( window, 'CSS', { value: () => ( {} ) } );
    Object.defineProperty( window, 'IS_SSR', {
        value: false,
        configurable: true,
        writable: true
    } );
    return DI;
}
