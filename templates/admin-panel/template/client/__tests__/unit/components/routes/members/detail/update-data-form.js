import * as directives from 'sham-ui-directives';
// eslint-disable-next-line max-len
import RoutesMembersDetailUpdateDataForm  from '../../../../../../src/components/routes/members/detail/update-data-form.sfc';
import renderer from 'sham-ui-test-helpers';

it( 'renders correctly', () => {
    const meta = renderer( RoutesMembersDetailUpdateDataForm, {}, {
        directives
    } );
    expect( meta.toJSON() ).toMatchSnapshot();
} );
