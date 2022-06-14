import setup, { app } from '../../../helpers';
import axios from 'axios';
jest.mock( 'axios' );

beforeEach( () => {
    jest.resetModules();
    jest.clearAllMocks();
} );

it( 'save name without change', async() => {
    expect.assertions( 3 );

    axios.useDefaultMocks();
    const DI = setup();

    history.pushState( {}, '', 'http://client.example.com/settings/' );
    await app.start( DI );
    app.click( '.panel.settings p:nth-of-type(1) .icon-pencil' );
    app.checkBody();

    axios
        .use( 'put', '/members/name', {
            'Status': 'OK',
            'Messages': [ axios.defaultMocksData.user.Name ]
        }, 200 );

    await app.form.submit();

    app.click( '[data-test-modal] [data-test-ok-button]' );
    await app.waitRendering();

    expect( axios.mocks.put ).toHaveBeenCalledTimes( 1 );
    app.checkBody();
} );
