import setup, { app } from '../../../helpers';
import axios from 'axios';
jest.mock( 'axios' );

beforeEach( () => {
    jest.resetModules();
    jest.clearAllMocks();
} );

it( 'success delete member', async() => {
    expect.assertions( 4 );

    const DI = setup();
    axios
        .useDefaultMocks()
        .use( 'get', '/validsession', {
            ...axios.defaultMocksData.user,
            IsSuperuser: true
        } )
        .use( 'get', 'admin/members', {
            members: [
                { ID: 1, Name: 'John Smith#1', Email: 'john.smith.1@test.com', IsSuperuser: false },
                { ID: 2, Name: 'John Smith#2', Email: 'john.smith.2@test.com', IsSuperuser: true },
                { ID: 3, Name: 'John Smith#3', Email: 'john.smith.3@test.com', IsSuperuser: false }
            ],
            meta: {
                total: 3,
                offset: 0,
                limit: 50
            }
        } );

    history.pushState( {}, '', 'http://client.example.com/members' );
    await app.start( DI );

    app.click( '[data-test-delete-button="1"]' );
    await app.waitRendering();

    app.checkMainPanel();

    axios
        .use( 'delete', 'admin/members/1' )
        .use( 'get', 'admin/members', {
            members: [
                { ID: 2, Name: 'John Smith#2', Email: 'john.smith.2@test.com', IsSuperuser: true },
                { ID: 3, Name: 'John Smith#3', Email: 'john.smith.3@test.com', IsSuperuser: false }
            ],
            meta: {
                total: 2,
                offset: 0,
                limit: 50
            }
        } );
    app.click( '[data-test-modal] [data-test-ok-button]' );
    await app.waitRendering();

    expect( axios.mocks.delete ).toHaveBeenCalledTimes( 1 );
    expect( axios.mocks.delete.mock.calls[ 0 ][ 0 ] ).toBe( 'admin/members/1' );
    app.checkMainPanel();
} );
