import setup, { app } from '../../../helpers';
import axios from 'axios';
jest.mock( 'axios' );

beforeEach( () => {
    jest.resetModules();
    jest.clearAllMocks();
} );

it( 'fail update tag data', async() => {
    expect.assertions( 1 );

    const DI = setup();
    axios
        .useDefaultMocks()
        .use( 'get', '/tags', {
            tags: [
                { 'id': '1', 'name': 'Hello world!', 'slug': 'hello-world' },
                { 'id': '6', 'name': 'Название тега', 'slug': 'nazvanie-tega' }
            ]
        } )
        .use( 'put', '/tags/6', {
            'Status': 'Bad request',
            'Messages': [ 'Can\'t update tag data' ]
        }, 400 );

    history.pushState( {}, '', 'http://client.example.com/articles/tags' );
    await app.start( DI );

    app.click( '[data-test-update-button="6"]' );
    await app.waitRendering();

    app.form.submit();

    app.click( '[data-test-modal] [data-test-ok-button]' );
    await app.waitRendering();

    app.checkMainPanel();
} );

it( 'fail update tag data (500 status code)', async() => {
    expect.assertions( 1 );

    const DI = setup();
    axios
        .useDefaultMocks()
        .use( 'get', '/tags', {
            tags: [
                { 'id': '1', 'name': 'Hello world!', 'slug': 'hello-world' },
                { 'id': '6', 'name': 'Название тега', 'slug': 'nazvanie-tega' }
            ]
        } )
        .use( 'put', '/tags/6', {}, 500 );

    history.pushState( {}, '', 'http://client.example.com/articles/tags' );
    await app.start( DI );

    app.click( '[data-test-update-button="6"]' );
    await app.waitRendering();

    app.form.submit();

    app.click( '[data-test-modal] [data-test-ok-button]' );
    await app.waitRendering();

    app.checkMainPanel();
} );
