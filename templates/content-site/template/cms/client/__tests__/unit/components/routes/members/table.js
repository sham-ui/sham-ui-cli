import RoutesMembersTable  from '../../../../../src/components/routes/members/table.sfc';
import renderer from 'sham-ui-test-helpers';

it( 'renders correctly', () => {
    const meta = renderer( RoutesMembersTable, {} );
    expect( meta.toJSON() ).toMatchSnapshot();
} );
