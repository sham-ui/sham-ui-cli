import setup, { app } from '../../../helpers';
import axios from 'axios';
jest.mock( 'axios' );
const DI = setup();

beforeEach( () => {
    jest.resetModules();
    jest.clearAllMocks();
} );

it( 'fail update member data', async() => {
    expect.assertions( 1 );

    axios
        .useDefaultMocks()
        .use( 'get', '/validsession', {
            ...axios.defaultMocksData.user,
            IsSuperuser: true
        } )
        .use( 'get', 'admin/members', {
            members: [
                { ID: 2, Name: 'John Smith#1', Email: 'john.smith.1@test.com', IsSuperuser: false }
            ],
            meta: {
                total: 1,
                offset: 0,
                limit: 50
            }
        } )
        .use( 'put', 'admin/members/2', {
            'Status': 'Bad request',
            'Messages': [ 'Can\'t update member data' ]
        }, 400 );

    history.pushState( {}, '', 'http://client.example.com/members' );
    await app.start( DI );

    app.click( '[data-test-update-button="2"]' );
    await app.waitRendering();

    app.form.submit();

    app.click( '[data-test-modal] [data-test-ok-button]' );
    await app.waitRendering();

    app.checkMainPanel();
} );

it( 'fail update member data (500 status code)', async() => {
    expect.assertions( 1 );

    axios
        .useDefaultMocks()
        .use( 'get', '/validsession', {
            ...axios.defaultMocksData.user,
            IsSuperuser: true
        } )
        .use( 'get', 'admin/members', {
            members: [
                { ID: 2, Name: 'John Smith#1', Email: 'john.smith.1@test.com', IsSuperuser: false }
            ],
            meta: {
                total: 1,
                offset: 0,
                limit: 50
            }
        } )
        .use( 'put', 'admin/members/2', {}, 500 );

    history.pushState( {}, '', 'http://client.example.com/members' );
    await app.start();

    app.click( '[data-test-update-button="2"]' );
    await app.waitRendering();

    app.form.submit();

    app.click( '[data-test-modal] [data-test-ok-button]' );
    await app.waitRendering();

    app.checkMainPanel();
} );
