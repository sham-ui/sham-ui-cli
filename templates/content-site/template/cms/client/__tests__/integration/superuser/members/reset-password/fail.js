import setup, { app } from '../../../helpers';
import axios from 'axios';
jest.mock( 'axios' );

beforeEach( () => {
    jest.resetModules();
    jest.clearAllMocks();
} );

it( 'fail reset member password', async() => {
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
                { ID: 2, Name: 'John Smith#1', Email: 'john.smith.1@test.com', IsSuperuser: false }
            ],
            meta: {
                total: 1,
                offset: 0,
                limit: 50
            }
        } )
        .use( 'put', 'admin/members/2/password', {
            'Status': 'Bad request',
            'Messages': [ 'Can\'t reset member password' ]
        }, 400 );

    history.pushState( {}, '', 'http://client.example.com/members' );
    await app.start( DI );

    app.click( '[data-test-update-button="2"]' );
    await app.waitRendering();

    app.form.fill( 'pass1', '' );
    app.form.fill( 'pass2', '' );
    app.click( '.form-layout:last-child [type="submit"]' );
    await app.waitRendering();

    app.click( '[data-test-modal] [data-test-ok-button]' );
    await app.waitRendering();

    app.checkMainPanel();
} );


it( 'fail reset member password (500 status code)', async() => {
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
        .use( 'put', 'admin/members/2/password', {}, 500 );

    history.pushState( {}, '', 'http://client.example.com/members' );
    await app.start();

    app.click( '[data-test-update-button="2"]' );
    await app.waitRendering();

    app.form.fill( 'pass1', '' );
    app.form.fill( 'pass2', '' );
    app.click( '.form-layout:last-child [type="submit"]' );
    await app.waitRendering();

    app.click( '[data-test-modal] [data-test-ok-button]' );
    await app.waitRendering();

    app.checkMainPanel();
} );
