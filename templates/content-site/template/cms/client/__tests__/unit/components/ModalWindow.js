import * as directives from 'sham-ui-directives';
import ModalWindow  from '../../../src/components/ModalWindow.sfc';
import renderer from 'sham-ui-test-helpers';

it( 'renders correctly', () => {
    const meta = renderer( ModalWindow, {
        directives
    } );
    expect( meta.toJSON() ).toMatchSnapshot();
} );

it( 'close', async() => {
    expect.assertions( 1 );
    const onClose = jest.fn();
    const meta = renderer( ModalWindow, {
        directives,
        onClose
    } );
    meta.component.container.querySelector( '[data-test-close-button]' ).click();
    await new Promise( resolve => setImmediate( resolve ) );
    expect( onClose ).toHaveBeenCalledTimes( 1 );
} );


it( 'cancel', async() => {
    expect.assertions( 1 );
    const onClose = jest.fn();
    const meta = renderer( ModalWindow, {
        directives,
        onClose
    } );
    meta.component.container.querySelector( '[data-test-cancel-button]' ).click();
    await new Promise( resolve => setImmediate( resolve ) );
    expect( onClose ).toHaveBeenCalledTimes( 1 );
} );


it( 'ok', async() => {
    expect.assertions( 1 );
    const onOk = jest.fn();
    const meta = renderer( ModalWindow, {
        directives,
        onOk
    } );
    meta.component.container.querySelector( '[data-test-ok-button]' ).click();
    await new Promise( resolve => setImmediate( resolve ) );
    expect( onOk ).toHaveBeenCalledTimes( 1 );
} );
