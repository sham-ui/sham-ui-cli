import Logo  from '../../../src/components/Logo.sht';
import renderer from 'sham-ui-test-helpers';

it( 'renders correctly', () => {
    const meta = renderer( Logo, {} );
    expect( meta.toJSON() ).toMatchSnapshot();
} );
