import setup, { app } from '../../../helpers';
import axios from 'axios';
jest.mock( 'axios' );

beforeEach( () => {
    jest.resetModules();
    jest.clearAllMocks();
} );

it( 'fail delete member', async() => {
    expect.assertions( 1 );

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
        } )
        .use( 'delete', 'admin/members/1', {
            'Status': 'Bad request',
            'Messages': [ 'Can\'t delete member' ]
        }, 400 );

    history.pushState( {}, '', 'http://client.example.com/members' );
    await app.start( DI );

    app.click( '[data-test-delete-button="1"]' );
    app.click( '[data-test-modal] [data-test-ok-button]' );
    await app.waitRendering();

    app.checkMainPanel();
} );

it( 'fail delete member (500 status code)', async() => {
    expect.assertions( 1 );

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
        } )
        .use( 'delete', 'admin/members/1', {}, 500 );

    history.pushState( {}, '', 'http://client.example.com/members' );
    await app.start();

    app.click( '[data-test-delete-button="1"]' );
    app.click( '[data-test-modal] [data-test-ok-button]' );
    await app.waitRendering();

    app.checkMainPanel();
} );
