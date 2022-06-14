import setup, { app } from '../helpers';
import axios from 'axios';
jest.mock( 'axios' );

beforeEach( () => {
    jest.resetModules();
    jest.clearAllMocks();
    window.matchMedia = jest.fn().mockImplementation(
        () => ( {
            addListener: jest.fn(),
            matches: true
        } )
    );
} );

afterEach( () => {
    delete window.matchMedia;
} );

it( 'logout', async() => {
    expect.assertions( 5 );

    axios
        .useDefaultMocks()
        .use( 'post', '/logout', '' );

    axios.useDefaultMocks();

    const DI = setup();

    await app.start( DI );

    axios.use( 'get', '/validsession', {}, 401 );
    app.click( '.icon-logout' );

    await app.waitRendering();
    app.checkBody();
    expect( window.location.href ).toBe( 'http://client.example.com/login' );
    expect( axios.mocks.post ).toHaveBeenCalledTimes( 1 );
    expect( axios.mocks.post.mock.calls[ 0 ] ).toHaveLength( 1 );
    expect( axios.mocks.post.mock.calls[ 0 ][ 0 ] ).toBe( '/logout' );
} );
