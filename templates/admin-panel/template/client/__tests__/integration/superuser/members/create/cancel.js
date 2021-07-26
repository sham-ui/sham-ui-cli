import setup, { app } from '../../../helpers';
import axios from 'axios';
jest.mock( 'axios' );

beforeEach( () => {
    jest.resetModules();
    jest.clearAllMocks();
} );

it( 'cancel create member', async() => {
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

    history.pushState( {}, '', 'http://client.example.com/members' );
    await app.start( DI );

    app.click( '[data-test-toggle-create-form]' );
    await app.waitRendering();

    app.form.fill( 'name', '' );
    app.form.fill( 'email', '' );
    app.form.fill( 'pass', '' );
    app.form.submit();

    app.click( '[data-test-modal] [data-test-cancel-button]' );
    await app.waitRendering();
    app.checkMainPanel();

    app.click( '[data-test-toggle-create-form]' );
    app.checkMainPanel();
} );
