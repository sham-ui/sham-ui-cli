import * as directives from 'sham-ui-directives';
// eslint-disable-next-line max-len
import RoutesMembersDetailUpdatePasswordForm  from '../../../../../../src/components/routes/members/detail/update-password-form.sfc';
import renderer from 'sham-ui-test-helpers';

it( 'renders correctly', () => {
    const meta = renderer( RoutesMembersDetailUpdatePasswordForm, {}, {
        directives
    } );
    expect( meta.toJSON() ).toMatchSnapshot();
} );
