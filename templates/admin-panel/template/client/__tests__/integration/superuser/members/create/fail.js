import setup, { app } from '../../../helpers';
import axios from 'axios';
jest.mock( 'axios' );

beforeEach( () => {
    jest.resetModules();
    jest.clearAllMocks();
} );

it( 'fail create member', async() => {
    expect.assertions( 1 );

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
        } )
        .use( 'post', 'admin/members', {
            'Status': 'Bad request',
            'Messages': [ 'Can\'t create member' ]
        }, 500 );

    history.pushState( {}, '', 'http://client.example.com/members' );
    await app.start( DI );

    app.click( '[data-test-toggle-create-form]' );
    await app.waitRendering();

    app.form.fill( 'name', 'John Smith' );
    app.form.fill( 'email', 'john.smith@test.com' );
    app.form.fill( 'pass', 'test' );
    app.form.submit();

    app.click( '[data-test-modal] [data-test-ok-button]' );
    await app.waitRendering();
    app.checkMainPanel();
} );

it( 'fail create member (500 status code)', async() => {
    expect.assertions( 1 );

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
        } )
        .use( 'post', 'admin/members', {}, 500 );

    history.pushState( {}, '', 'http://client.example.com/members' );
    await app.start( DI );

    app.click( '[data-test-toggle-create-form]' );
    await app.waitRendering();

    app.form.fill( 'name', 'John Smith' );
    app.form.fill( 'email', 'john.smith@test.com' );
    app.form.fill( 'pass', 'test' );
    app.form.submit();

    app.click( '[data-test-modal] [data-test-ok-button]' );
    await app.waitRendering();
    app.checkMainPanel();
} );
