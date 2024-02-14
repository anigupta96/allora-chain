package emissions

import (
	"errors"
	"fmt"
)

var ErrIntegerUnderflowDelegator = errors.New(Err_ErrIntegerUnderflowDelegator.String())
var ErrIntegerUnderflowBonds = errors.New(Err_ErrIntegerUnderflowBonds.String())
var ErrIntegerUnderflowTarget = errors.New(Err_ErrIntegerUnderflowTarget.String())
var ErrIntegerUnderflowTopicStake = errors.New(Err_ErrIntegerUnderflowTopicStake.String())
var ErrIntegerUnderflowTotalStake = errors.New(Err_ErrIntegerUnderflowTotalStake.String())
var ErrIterationLengthDoesNotMatch = errors.New(Err_ErrIterationLengthDoesNotMatch.String())
var ErrInvalidTopicId = fmt.Errorf(Err_ErrInvalidTopicId.String())
var ErrReputerAlreadyRegistered = fmt.Errorf(Err_ErrReputerAlreadyRegistered.String())
var ErrWorkerAlreadyRegistered = fmt.Errorf(Err_ErrWorkerAlreadyRegistered.String())
var ErrInsufficientStakeToRegister = fmt.Errorf(Err_ErrInsufficientStakeToRegister.String())
var ErrLibP2PKeyRequired = fmt.Errorf(Err_ErrLibP2PKeyRequired.String())
var ErrAddressNotRegistered = fmt.Errorf(Err_ErrAddressNotRegistered.String())
var ErrStakeTargetNotRegistered = fmt.Errorf(Err_ErrStakeTargetNotRegistered.String())
var ErrTopicIdOfStakerAndTargetDoNotMatch = fmt.Errorf(Err_ErrInvalidTopicId.String())
var ErrInsufficientStakeToRemove = fmt.Errorf(Err_ErrInsufficientStakeToRemove.String())
var ErrNoStakeToRemove = fmt.Errorf(Err_ErrNoStakeToRemove.String())
var ErrDoNotSetMapValueToZero = fmt.Errorf(Err_ErrDoNotSetMapValueToZero.String())
var ErrBlockHeightNegative = fmt.Errorf(Err_ErrBlockHeightNegative.String())
var ErrBlockHeightLessThanPrevious = fmt.Errorf(Err_ErrBlockHeightLessThanPrevious.String())
var ErrModifyStakeBeforeBondLessThanAmountModified = fmt.Errorf(Err_ErrModifyStakeBeforeBondLessThanAmountModified.String())
var ErrModifyStakeBeforeSumGreaterThanSenderStake = fmt.Errorf(Err_ErrModifyStakeBeforeSumGreaterThanSenderStake.String())
var ErrModifyStakeSumBeforeNotEqualToSumAfter = fmt.Errorf(Err_ErrModifyStakeSumBeforeNotEqualToSumAfter.String())
var ErrScalarMultiplyNegative = fmt.Errorf(Err_ErrScalarMultiplyNegative.String())
var ErrDivideMapValuesByZero = fmt.Errorf(Err_ErrDivideMapValuesByZero.String())