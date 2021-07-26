import setup, { app } from '../../../helpers';
import axios from 'axios';
jest.mock( 'axios' );

beforeEach( () => {
    jest.resetModules();
    jest.clearAllMocks();
} );

it( 'cancel edit password', async() => {
    expect.assertions( 3 );

    axios.useDefaultMocks();
    const DI = setup();

    history.pushState( {}, '', 'http://client.example.com/settings/' );
    await app.start( DI );
    app.click( '.panel.settings p:nth-of-type(3) .icon-pencil' );
    app.checkBody();

    app.form.fill( 'pass1', '' );
    app.form.fill( 'pass2', '' );
    await app.form.submit();

    app.click( '[data-test-modal] [data-test-cancel-button]' );

    await app.waitRendering();

    expect( axios.mocks.put ).toHaveBeenCalledTimes( 0 );
    app.checkBody();
} );
