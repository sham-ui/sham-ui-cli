import setup, { app } from '../../helpers';
import axios from 'axios';
jest.mock( 'axios' );

beforeEach( () => {
    jest.resetModules();
    jest.clearAllMocks();
} );

it( 'can go from home', async() => {
    expect.assertions( 2 );

    const DI = setup();
    axios
        .useDefaultMocks()
        .use( 'get', '/categories', {
            categories: []
        } );

    await app.start( DI );
    app.click( '.sideleft .icon-flow-split' );
    await app.waitRendering();
    app.checkBody();
    expect( window.location.href ).toBe( 'http://client.example.com/articles/categories' );
} );
