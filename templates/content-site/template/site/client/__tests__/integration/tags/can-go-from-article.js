import setupUnsafe from 'sham-ui-unsafe';
import axios from 'axios';
import setup, { app } from '../helpers';
jest.mock( 'axios' );

beforeEach( () => {
    jest.resetModules();
    jest.clearAllMocks();
} );

it( 'go from article', async() => {
    expect.assertions( 4 );

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
        .use( 'get', '/articles', {
            articles: [
                {
                    'title': 'Первая тестовая статья',
                    'slug': 'pervya-testovaya-statya',
                    'category': {
                        'name': 'Первая',
                        'slug': 'pervya'
                    },
                    'content': 'Short content',
                    'createdAt': '2022-05-24 19:54:20 +0000 UTC'
                }
            ],
            meta: {
                total: 100,
                limit: 20,
                tag: 'Быт'
            }
        } )
    ;

    history.pushState( {}, '', 'http://client.example.com/pervya-testovaya-statya' );

    await app.start( DI );

    app.click( 'a[data-test-tag="byt"]' );
    await app.waitRendering();

    expect( window.location.href ).toBe( 'http://client.example.com/tag/byt/1/' );
    expect( axios.mocks.get ).toHaveBeenCalledTimes( 2 );
    expect( axios.mocks.get.mock.calls[ 1 ] ).toEqual( [
        '/articles',
        { 'limit': 9, 'offset': 0, 'tag': 'byt' }
    ] );

    app.checkBody();
} );

it( 'go from article (ssr & rehydrate)', async() => {
    expect.assertions( 4 );

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
        .use( 'get', '/articles', {
            articles: [
                {
                    'title': 'Первая тестовая статья',
                    'slug': 'pervya-testovaya-statya',
                    'category': {
                        'name': 'Первая',
                        'slug': 'pervya'
                    },
                    'content': 'Short content',
                    'createdAt': '2022-05-24 19:54:20 +0000 UTC'
                }
            ],
            meta: {
                total: 100,
                limit: 20,
                tag: 'Быт'
            }
        } )
    ;

    await app.ssrAndRehydrate( 'http://client.example.com/pervya-testovaya-statya' );

    app.click( 'a[data-test-tag="byt"]' );
    await app.waitRendering();

    expect( window.location.href ).toBe( 'http://client.example.com/tag/byt/1/' );
    expect( axios.mocks.get ).toHaveBeenCalledTimes( 2 );
    expect( axios.mocks.get.mock.calls[ 1 ] ).toEqual( [
        '/articles',
        { 'limit': 9, 'offset': 0, 'tag': 'byt' }
    ] );

    app.checkBody();
} );
