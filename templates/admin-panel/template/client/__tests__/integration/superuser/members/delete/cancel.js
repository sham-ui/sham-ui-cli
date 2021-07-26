import setup, { app } from '../../../helpers';
import axios from 'axios';
jest.mock( 'axios' );

beforeEach( () => {
    jest.resetModules();
    jest.clearAllMocks();
} );

it( 'cancel delete member', async() => {
    expect.assertions( 2 );

    const DI = setup();
    axios
        .useDefaultMocks()
        .use( 'get', '/validsession', {
            ...axios.defaultMocksData.user,
            IsSuperuser: true
        } )
        .use( 'get', 'admin/members', {
            members: [
                { ID: 1, Name: 'John Smith#1', Email: 'john.smith.1@test.com', IsSuperuser: false }
            ],
            meta: {
                total: 1,
                offset: 0,
                limit: 50
            }
        } );

    history.pushState( {}, '', 'http://client.example.com/members' );
    await app.start( DI );

    app.click( '[data-test-delete-button="1"]' );
    app.click( '[data-test-modal] [data-test-cancel-button]' );
    await app.waitRendering();

    expect( axios.mocks.delete ).toHaveBeenCalledTimes( 0 );
    app.checkMainPanel();
} );
