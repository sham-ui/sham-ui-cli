import * as directives from 'sham-ui-directives';
// eslint-disable-next-line max-len
import RoutesSettingsFormPassword  from '../../../../../../src/components/routes/settings/form/password.sfc';
import renderer from 'sham-ui-test-helpers';
import { createDI } from 'sham-ui';

it( 'renders correctly', () => {
    const DI = createDI();
    const meta = renderer( RoutesSettingsFormPassword, {}, {
        DI,
        directives
    } );
    expect( meta.toJSON() ).toMatchSnapshot();
} );

it( 'display errors', async() => {
    expect.assertions( 2 );

    const DI = createDI();
    const updateMock = jest.fn();
    DI.bind( 'store', {
        updateMemberPassword: updateMock.mockReturnValueOnce( Promise.reject( {} ) )
    } );

    const meta = renderer( RoutesSettingsFormPassword, {}, {
        DI,
        directives
    } );

    const formData = {
        pass1: 'admin1@gmail.com',
        pass2: 'admin1@gmail.com'
    };
    const { ctx } = meta;
    ctx.container.querySelector( '[name="pass1"]' ).value = formData.pass1;
    ctx.container.querySelector( '[name="pass2"]' ).value = formData.pass2;
    ctx.container.querySelector( '[type="submit"]' ).click();
    ctx.container.querySelector( '[data-test-modal] [data-test-ok-button]' ).click();

    await new Promise( resolve => setImmediate( resolve ) );

    expect( updateMock.mock.calls ).toHaveLength( 1 );
    expect( meta.toJSON() ).toMatchSnapshot();
} );
