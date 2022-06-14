import setupUnsafe from 'sham-ui-unsafe';
import axios from 'axios';
import setup, { app, setupSSR } from '../helpers';
jest.mock( 'axios' );

beforeEach( () => {
    jest.resetModules();
    jest.clearAllMocks();
} );

it( 'display', async() => {
    expect.assertions( 4 );

    const DI = setup();
    setupUnsafe( DI );

    axios
        .useDefaultMocks()
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

    history.pushState( {}, '', 'http://client.example.com/tag/byt/1' );

    await app.start( DI );
    app.checkBody();
    expect( window.location.href ).toBe( 'http://client.example.com/tag/byt/1' );
    expect( axios.mocks.get ).toHaveBeenCalledTimes( 1 );
    expect( axios.mocks.get.mock.calls[ 0 ] ).toEqual( [
        '/articles',
        { 'limit': 9, 'offset': 0, 'tag': 'byt' }
    ] );
} );


it( 'ssr', async() => {
    expect.assertions( 3 );
    axios
        .useDefaultMocks()
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
    setupSSR();
    await app.ssrAndCheck( 'http://client.example.com/tag/byt/1' );
    expect( axios.mocks.get ).toHaveBeenCalledTimes( 1 );
    expect( axios.mocks.get.mock.calls[ 0 ] ).toEqual( [
        '/articles',
        { 'limit': 9, 'offset': 0, 'tag': 'byt' }
    ] );
} );


it( 'ssr & rehydrate', async() => {
    expect.assertions( 4 );

    axios
        .useDefaultMocks()
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

    await app.ssrAndRehydrate( 'http://client.example.com/tag/byt/1' );

    app.checkBody();
    expect( window.location.href ).toBe( 'http://client.example.com/tag/byt/1' );
    expect( axios.mocks.get ).toHaveBeenCalledTimes( 1 );
    expect( axios.mocks.get.mock.calls[ 0 ] ).toEqual( [
        '/articles',
        { 'limit': 9, 'offset': 0, 'tag': 'byt' }
    ] );
} );
