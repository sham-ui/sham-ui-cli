import setup, { app } from '../../../helpers';
import axios from 'axios';
jest.mock( 'axios' );

beforeEach( () => {
    jest.resetModules();
    jest.clearAllMocks();
} );

it( 'fail create category', async() => {
    expect.assertions( 1 );

    const DI = setup();
    axios
        .useDefaultMocks()
        .use( 'get', '/categories', {
            categories: []
        } )
        .use( 'post', 'categories', {
            'Status': 'Bad request',
            'Messages': [ 'Can\'t create category' ]
        }, 500 );

    history.pushState( {}, '', 'http://client.example.com/articles/categories' );
    await app.start( DI );

    app.click( '[data-test-toggle-create-form]' );
    await app.waitRendering();

    app.form.fill( 'name', 'travel' );
    app.form.submit();

    app.click( '[data-test-modal] [data-test-ok-button]' );
    await app.waitRendering();
    app.checkMainPanel();
} );

it( 'fail create category (500 status code)', async() => {
    expect.assertions( 1 );

    const DI = setup();
    axios
        .useDefaultMocks().
        use( 'get', '/categories', {
            categories: []
        } )
        .use( 'post', 'categories', {}, 500 );

    history.pushState( {}, '', 'http://client.example.com/articles/categories' );
    await app.start( DI );

    app.click( '[data-test-toggle-create-form]' );
    await app.waitRendering();

    app.form.fill( 'name', 'travel' );
    app.form.submit();

    app.click( '[data-test-modal] [data-test-ok-button]' );
    await app.waitRendering();
    app.checkMainPanel();
} );
