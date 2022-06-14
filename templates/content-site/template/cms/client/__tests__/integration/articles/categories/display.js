import setup, { app } from '../../helpers';
import axios from 'axios';
jest.mock( 'axios' );

beforeEach( () => {
    jest.resetModules();
    jest.clearAllMocks();
} );

it( 'display page', async() => {
    expect.assertions( 2 );

    const DI = setup();
    axios
        .useDefaultMocks()
        .use( 'get', '/categories', {
            categories: [
                { 'id': '4', 'name': 'Hello world!', 'slug': 'hello-world' },
                { 'id': '6', 'name': 'Название категории', 'slug': 'nazvanie-kategorii' }
            ]
        } );

    history.pushState( {}, '', 'http://client.example.com/articles/categories' );
    await app.start( DI );
    app.checkMainPanel();
    expect( window.location.href ).toBe( 'http://client.example.com/articles/categories' );
} );
