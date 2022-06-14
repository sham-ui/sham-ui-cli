import setup, { app } from '../../../helpers';
import axios from 'axios';
jest.mock( 'axios' );

beforeEach( () => {
    jest.resetModules();
    jest.clearAllMocks();
} );

it( 'fail delete tag', async() => {
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
        .use( 'delete', '/tags/1', {
            'Status': 'Bad request',
            'Messages': [ 'Can\'t delete tag' ]
        }, 400 );

    history.pushState( {}, '', 'http://client.example.com/articles/tags' );
    await app.start( DI );

    app.click( '[data-test-delete-button="1"]' );
    app.click( '[data-test-modal] [data-test-ok-button]' );
    await app.waitRendering();

    app.checkMainPanel();
} );

it( 'fail delete tag (500 status code)', async() => {
    expect.assertions( 1 );

    axios
        .useDefaultMocks()
        .use( 'get', '/tags', {
            tags: [
                { 'id': '1', 'name': 'Hello world!', 'slug': 'hello-world' },
                { 'id': '6', 'name': 'Название тэга', 'slug': 'nazvanie-tega' }
            ]
        } )
        .use( 'delete', '/tags/1', {}, 500 );

    history.pushState( {}, '', 'http://client.example.com/articles/tags' );
    await app.start();

    app.click( '[data-test-delete-button="1"]' );
    app.click( '[data-test-modal] [data-test-ok-button]' );
    await app.waitRendering();

    app.checkMainPanel();
} );
