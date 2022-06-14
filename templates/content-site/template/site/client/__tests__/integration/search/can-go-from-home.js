import setupUnsafe from 'sham-ui-unsafe';
import axios from 'axios';
import setup, { app } from '../helpers';
jest.mock( 'axios' );

beforeEach( () => {
    jest.resetModules();
    jest.clearAllMocks();
} );

it( 'go from home', async() => {
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
                limit: 9,
                category: 'Первая'
            }
        } )
    ;

    history.pushState( {}, '', 'http://client.example.com/' );

    await app.start( DI );

    app.click( '.search-icon' );
    app.form.fill( 'query', 'content' );
    await app.form.submit();

    expect( window.location.href ).toBe( 'http://client.example.com/search/1/?q=content' );
    expect( axios.mocks.get ).toHaveBeenCalledTimes( 2 );
    expect( axios.mocks.get.mock.calls[ 1 ] ).toEqual( [
        '/articles',
        { 'limit': 9, 'offset': 0, q: 'content' }
    ] );
    app.checkBody();
} );

it( 'go from home (ssr & rehydrate)', async() => {
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
                limit: 9,
                category: 'Первая'
            }
        } )
    ;

    await app.ssrAndRehydrate( 'http://client.example.com/' );

    app.click( '.search-icon' );
    app.form.fill( 'query', 'content' );
    await app.form.submit();

    expect( window.location.href ).toBe( 'http://client.example.com/search/1/?q=content' );
    expect( axios.mocks.get ).toHaveBeenCalledTimes( 2 );
    expect( axios.mocks.get.mock.calls[ 1 ] ).toEqual( [
        '/articles',
        { 'limit': 9, 'offset': 0, q: 'content' }
    ] );
    app.checkBody();
} );
