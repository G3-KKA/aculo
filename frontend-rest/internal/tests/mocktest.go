package testutils

//go:generate mockery --name=Testutils --dir=. --outpkg=mmk --filename=mock_testutils.go --output=./mocks/testutils --structname MockTestutils --inpackage=true
type Testutils interface {
}

//go:generate mockery --name=Interface --dir=. --structname MockInterface --filename=mock_interface.go --inpackage=true
type Interface interface {
}
