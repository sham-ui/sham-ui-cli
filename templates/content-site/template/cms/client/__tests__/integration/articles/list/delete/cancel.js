import setup, { app } from '../../../helpers';
import axios from 'axios';
jest.mock( 'axios' );

beforeEach( () => {
    jest.resetModules();
    jest.clearAllMocks();
} );

it( 'cancel delete article', async() => {
    expect.assertions( 2 );

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
                }
            ],
            'meta': {
                'limit': 50, 'offset': 0, 'total': 1
            }
        } );

    history.pushState( {}, '', 'http://client.example.com/articles' );
    await app.start( DI );

    app.click( '[data-test-delete-button="10"]' );
    app.click( '[data-test-modal] [data-test-cancel-button]' );
    await app.waitRendering();

    expect( axios.mocks.delete ).toHaveBeenCalledTimes( 0 );
    app.checkMainPanel();
} );
