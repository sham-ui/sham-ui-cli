import { start, createDI } from 'sham-ui';
import pretty from 'pretty';
import initializer from '../../src/initializers/main';

export const app = {
    async start( DI, waitRendering = true ) {
        if ( !DI ) {
            DI = createDI();
        }
        initializer( DI );
        start( DI );
        if ( waitRendering ) {
            await this.waitRendering();
        }
    },
    async waitRendering() {
        await new Promise( resolve => setImmediate( resolve ) );
    },
    click( selector ) {
        document.querySelector( selector ).click();
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

function resetStorage( DI ) {
    const storage = DI.resolve( 'session:storage' );
    if ( storage ) {
        storage.reset();
    }
}

function setupRouter() {
    delete window.__NAVIGO_WINDOW_LOCATION_MOCK__;
    history.pushState( {}, '', '' );
}

export default function( DI ) {
    if ( !DI ) {
        DI = createDI();
    }
    setupRAF();
    clearBody();
    resetStorage( DI );
    setupRouter();
    Object.defineProperty( window, 'CSS', { value: () => ( {} ) } );
    return DI;
}
