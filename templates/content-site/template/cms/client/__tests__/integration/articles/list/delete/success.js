import setup, { app } from '../../../helpers';
import axios from 'axios';
jest.mock( 'axios' );

beforeEach( () => {
    jest.resetModules();
    jest.clearAllMocks();
} );

it( 'success delete article', async() => {
    expect.assertions( 4 );

    const DI = setup();
    axios
        .useDefaultMocks()
        .use( 'get', '/categories', {
            'categories': [
                { 'id': '1', 'name': 'Кухня', 'slug': 'kukhnia' },
                { 'id': '2', 'name': 'Быт', 'slug': 'byt' }
            ]
        } )
        .use( 'get', '/articles', {
            'articles': [
                {
                    'id': '10',
                    'title': 'Тест2',
                    'slug': 'test2',
                    'category_id': '1',
                    'published_at': '2022-04-25T19:34:23.619+07:00'
                },
                {
                    'id': '2',
                    'title': 'Вторая',
                    'slug': 'second',
                    'category_id': '2',
                    'published_at': '2022-03-22T08:26:23.619+07:00'
                }
            ],
            'meta': {
                'limit': 50, 'offset': 0, 'total': 2
            }
        } );

    history.pushState( {}, '', 'http://client.example.com/articles' );
    await app.start( DI );

    app.click( '[data-test-delete-button="10"]' );
    await app.waitRendering();

    app.checkMainPanel();

    axios
        .use( 'delete', '/articles/10' )
        .use( 'get', '/articles', {
            'articles': [
                {
                    'id': '2',
                    'title': 'Вторая',
                    'slug': 'second',
                    'category_id': '2',
                    'published_at': '2022-03-22T08:26:23.619+07:00'
                }
            ],
            'meta': {
                'limit': 50, 'offset': 0, 'total': 1
            }
        } );
    app.click( '[data-test-modal] [data-test-ok-button]' );
    await app.waitRendering();

    expect( axios.mocks.delete ).toHaveBeenCalledTimes( 1 );
    expect( axios.mocks.delete.mock.calls[ 0 ][ 0 ] ).toBe( '/articles/10' );
    app.checkMainPanel();
} );
