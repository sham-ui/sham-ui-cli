import setup, { app } from '../helpers';
import axios from 'axios';
jest.mock( 'axios' );

beforeEach( () => {
    jest.resetModules();
    jest.clearAllMocks();
} );

it( 'display', async() => {
    expect.assertions( 2 );

    const DI = setup();
    axios.useDefaultMocks();

    history.pushState( {}, '', 'http://client.example.com/settings/' );
    await app.start( DI );
    app.checkBody();
    expect( window.location.href ).toBe( 'http://client.example.com/settings/' );
} );
