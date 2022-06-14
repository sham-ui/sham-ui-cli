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
        .use( 'get', '/validsession', {
            ...axios.defaultMocksData.user,
            IsSuperuser: true
        } )
        .use( 'get', 'admin/members', {
            members: [],
            meta: {
                total: 0,
                offset: 0,
                limit: 50
            }
        } );

    await app.start( DI );
    app.click( '.sideleft .icon-users' );
    await app.waitRendering();
    app.checkBody();
    expect( window.location.href ).toBe( 'http://client.example.com/members' );
} );
