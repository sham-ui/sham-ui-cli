import setup, { app } from '../../helpers';
import axios from 'axios';
jest.mock( 'axios' );

beforeEach( () => {
    jest.resetModules();
    jest.clearAllMocks();
} );

it( 'pagination work', async() => {
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
                { ID: 1, Name: 'John Smith#1', Email: 'john.smith.1@test.com', IsSuperuser: false },
                { ID: 2, Name: 'John Smith#2', Email: 'john.smith.2@test.com', IsSuperuser: true },
                { ID: 3, Name: 'John Smith#3', Email: 'john.smith.3@test.com', IsSuperuser: false }
            ],
            meta: {
                total: 51,
                offset: 0,
                limit: 50
            }
        } );
    history.pushState( {}, '', 'http://client.example.com/members' );
    await app.start( DI );

    app.click( '[data-test-pagination-page="1"]' );
    app.checkMainPanel();

    axios
        .use( 'get', 'admin/members', {
            members: [
                { ID: 4, Name: 'John Smith#4', Email: 'john.smith.4@test.com', IsSuperuser: false },
                { ID: 5, Name: 'John Smith#5', Email: 'john.smith.5@test.com', IsSuperuser: true }
            ],
            meta: {
                total: 51,
                offset: 50,
                limit: 50
            }
        } );
    app.click( '[data-test-pagination-page="2"]' );
    await app.waitRendering();
    app.checkMainPanel();
} );
