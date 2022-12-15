import * as directives from 'sham-ui-directives';
// eslint-disable-next-line max-len
import RoutesSettingsFormEmail  from '../../../../../../src/components/routes/settings/form/email.sfc';
import renderer from 'sham-ui-test-helpers';
import { createDI } from 'sham-ui';

it( 'renders correctly', () => {
    const DI = createDI();
    const meta = renderer( RoutesSettingsFormEmail, {}, {
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
        updateMemberEmail: updateMock.mockReturnValueOnce( Promise.reject( {} ) )
    } );

    const meta = renderer( RoutesSettingsFormEmail, {}, {
        DI,
        directives
    } );

    const formData = {
        email1: 'admin1@gmail.com',
        email2: 'admin1@gmail.com'
    };
    const { ctx } = meta;
    ctx.container.querySelector( '[name="email1"]' ).value = formData.email1;
    ctx.container.querySelector( '[name="email2"]' ).value = formData.email2;
    ctx.container.querySelector( '[type="submit"]' ).click();
    ctx.container.querySelector( '[data-test-modal] [data-test-ok-button]' ).click();

    await new Promise( resolve => setImmediate( resolve ) );

    expect( updateMock.mock.calls ).toHaveLength( 1 );
    expect( meta.toJSON() ).toMatchSnapshot();
} );
