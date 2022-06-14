import setup, { app } from '../helpers';
import axios from 'axios';
jest.mock( 'axios' );

beforeEach( () => {
    jest.resetModules();
    jest.clearAllMocks();
} );

it( 'authenticated member can\'t visit login page', async() => {
    expect.assertions( 2 );

    axios.useDefaultMocks();

    const DI = setup();

    history.pushState( {}, '', 'http://client.example.com/login/' );
    await app.start( DI );
    app.checkBody();
    expect( window.location.href ).toBe( 'http://client.example.com/' );
} );
