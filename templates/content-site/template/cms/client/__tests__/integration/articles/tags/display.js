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
        .use( 'get', '/tags', {
            tags: [
                { 'id': '4', 'name': 'Hello world!', 'slug': 'hello-world' },
                { 'id': '6', 'name': 'Название тега', 'slug': 'nazvanie-tega' }
            ]
        } );

    history.pushState( {}, '', 'http://client.example.com/articles/tags' );
    await app.start( DI );
    app.checkMainPanel();
    expect( window.location.href ).toBe( 'http://client.example.com/articles/tags' );
} );
