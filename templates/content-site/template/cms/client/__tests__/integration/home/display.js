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

it( 'display', async() => {
    expect.assertions( 2 );

    axios.useDefaultMocks();
    const DI = setup();

    await app.start( DI );
    app.checkBody();
    expect( window.location.href ).toBe( 'http://client.example.com/' );
} );
