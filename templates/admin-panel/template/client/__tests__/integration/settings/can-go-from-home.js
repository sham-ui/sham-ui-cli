import setup, { app } from '../helpers';
import axios from 'axios';
jest.mock( 'axios' );

beforeEach( () => {
    jest.resetModules();
    jest.clearAllMocks();
} );

it( 'can go from home', async() => {
    expect.assertions( 2 );

    const DI = setup();
    axios.useDefaultMocks();

    await app.start( DI );
    app.click( '.link-profile' );
    app.click( '.icon-cog' );
    await app.waitRendering();
    app.checkBody();
    expect( window.location.href ).toBe( 'http://client.example.com/settings' );
} );
