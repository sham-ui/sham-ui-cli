import setup, { app } from '../../../helpers';
import axios from 'axios';
jest.mock( 'axios' );

beforeEach( () => {
    jest.resetModules();
    jest.clearAllMocks();
} );

it( 'cancel create category', async() => {
    expect.assertions( 2 );

    const DI = setup();
    axios
        .useDefaultMocks()
        .use( 'get', '/categories', {
            categories: []
        } );

    history.pushState( {}, '', 'http://client.example.com/articles/categories' );
    await app.start( DI );

    app.click( '[data-test-toggle-create-form]' );
    await app.waitRendering();

    app.form.fill( 'name', '' );
    app.form.submit();

    app.click( '[data-test-modal] [data-test-cancel-button]' );
    await app.waitRendering();
    app.checkMainPanel();

    app.click( '[data-test-toggle-create-form]' );
    app.checkMainPanel();
} );
