import setupUnsafe from 'sham-ui-unsafe';
import axios from 'axios';
import setup, { app } from '../../helpers';
jest.mock( 'axios' );

beforeEach( () => {
    jest.resetModules();
    jest.clearAllMocks();
} );

it( 'go to page', async() => {
    expect.assertions( 7 );

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

    history.pushState( {}, '', 'http://client.example.com/category/pervya/5' );

    await app.start( DI );

    expect( window.location.href ).toBe( 'http://client.example.com/category/pervya/5' );
    expect( axios.mocks.get ).toHaveBeenCalledTimes( 1 );
    expect( axios.mocks.get.mock.calls[ 0 ] ).toEqual( [
        '/articles',
        { 'limit': 9, 'offset': 36, 'category': 'pervya' }
    ] );

    app.click( '[data-test-page="6"]' );
    await app.waitRendering();

    app.checkBody();
    expect( axios.mocks.get ).toHaveBeenCalledTimes( 2 );
    expect( axios.mocks.get.mock.calls[ 1 ] ).toEqual( [
        '/articles',
        { 'limit': 9, 'offset': 45, 'category': 'pervya' }
    ] );
    expect( window.location.href ).toBe( 'http://client.example.com/category/pervya/6/' );
} );

it( 'go to page (ssr & rehydrate)', async() => {
    expect.assertions( 7 );

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

    await app.ssrAndRehydrate( 'http://client.example.com/category/pervya/5' );

    expect( window.location.href ).toBe( 'http://client.example.com/category/pervya/5' );
    expect( axios.mocks.get ).toHaveBeenCalledTimes( 1 );
    expect( axios.mocks.get.mock.calls[ 0 ] ).toEqual( [
        '/articles',
        { 'limit': 9, 'offset': 36, 'category': 'pervya' }
    ] );

    app.click( '[data-test-page="6"]' );
    await app.waitRendering();

    app.checkBody();
    expect( axios.mocks.get ).toHaveBeenCalledTimes( 2 );
    expect( axios.mocks.get.mock.calls[ 1 ] ).toEqual( [
        '/articles',
        { 'limit': 9, 'offset': 45, 'category': 'pervya' }
    ] );
    expect( window.location.href ).toBe( 'http://client.example.com/category/pervya/6/' );
} );
