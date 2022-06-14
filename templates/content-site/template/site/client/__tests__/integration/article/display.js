import setupUnsafe from 'sham-ui-unsafe';
import axios from 'axios';
import setup, { app, setupSSR } from '../helpers';
jest.mock( 'axios' );

beforeEach( () => {
    jest.resetModules();
    jest.clearAllMocks();
} );

it( 'display', async() => {
    expect.assertions( 2 );

    const DI = setup();
    setupUnsafe( DI );

    axios
        .useDefaultMocks()
        .use( 'get', '/articles/pervya-testovaya-statya', {
            'title': 'Первая тестовая статья',
            'slug': 'pervya-testovaya-statya',
            'category': {
                'name': 'Первая',
                'slug': 'pervya'
            },
            'tags': [ {
                'name': 'Быт',
                'slug': 'byt'
            } ],
            'content': '<p>Short <b>content</b></p>',
            'createdAt': '2022-05-24 19:54:20 +0000 UTC'
        } )
    ;

    history.pushState( {}, '', 'http://client.example.com/pervya-testovaya-statya' );

    await app.start( DI );
    app.checkBody();
    expect( window.location.href ).toBe( 'http://client.example.com/pervya-testovaya-statya' );
} );


it( 'ssr', async() => {
    expect.assertions( 1 );
    axios
        .useDefaultMocks()
        .use( 'get', '/articles/pervya-testovaya-statya', {
            'title': 'Первая тестовая статья',
            'slug': 'pervya-testovaya-statya',
            'category': {
                'name': 'Первая',
                'slug': 'pervya'
            },
            'tags': [ {
                'name': 'Быт',
                'slug': 'byt'
            } ],
            'content': '<p>Short <b>content</b></p>',
            'createdAt': '2022-05-24 19:54:20 +0000 UTC'
        } )
    ;
    setupSSR();
    await app.ssrAndCheck( 'http://client.example.com/pervya-testovaya-statya' );
} );


it( 'ssr & rehydrate', async() => {
    expect.assertions( 2 );

    axios
        .useDefaultMocks()
        .use( 'get', '/articles/pervya-testovaya-statya', {
            'title': 'Первая тестовая статья',
            'slug': 'pervya-testovaya-statya',
            'category': {
                'name': 'Первая',
                'slug': 'pervya'
            },
            'tags': [ {
                'name': 'Быт',
                'slug': 'byt'
            } ],
            'content': '<p>Short <b>content</b></p>',
            'createdAt': '2022-05-24 19:54:20 +0000 UTC'
        } )
    ;

    await app.ssrAndRehydrate( 'http://client.example.com/pervya-testovaya-statya' );

    app.checkBody();
    expect( window.location.href ).toBe( 'http://client.example.com/pervya-testovaya-statya' );
} );
