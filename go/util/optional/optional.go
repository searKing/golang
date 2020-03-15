// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package optional

import (
	"errors"

	"github.com/searKing/golang/go/util/function/consumer"
	"github.com/searKing/golang/go/util/function/predicate"
	"github.com/searKing/golang/go/util/object"
)

var (
	ErrorNoValuePresent = errors.New("no value present")
)

// Optional is a container object which may or may not contain a non-{@code null} value.
// If a value is present, {@code isPresent()} returns {@code true}. If no
// value is present, the object is considered <i>empty</i> and
// {@code isPresent()} returns {@code false}.
type Optional struct {
	isPresent bool
	value     interface{}
}

// Empty returns an empty {@code Optional} instance.  No value is present for this
// {@code Optional}.
func Empty() Optional {
	return Optional{}
}

// Of Returns an {@code Optional} describing the given value.
func Of(value interface{}) Optional {
	return Optional{
		isPresent: true,
		value:     value,
	}
}

// Of Returns an {@code Optional} describing the given value, if
// non-{@code null}, otherwise returns an empty {@code Optional}.
func OfNillable(value interface{}) Optional {
	if value == nil {
		return Empty()
	}
	return Of(value)
}

// Get returns the value if a value is present, otherwise throws
// {@code ErrorNoValuePresent}.
func (o Optional) Get() interface{} {
	return o.value
}

// IsPresent returns {@code true} if a value is present, otherwise {@code false}.
func (o Optional) IsPresent() bool {
	return o.isPresent
}

type EmptyConsumer interface {
	Run()
}

/**
 * If a value is present, performs the given action with the value,
 * otherwise does nothing.
 *
 * @param action the action to be performed, if a value is present
 * @throws ErrorNilPointer if value is present and the given action is
 *         {@code null}
 */
func (o Optional) IfPresent(action consumer.Consumer) {
	if o.IsPresent() {
		action.Accept(o.value)
	}
}

/**
 * If a value is present, performs the given action with the value,
 * otherwise performs the given empty-based action.
 *
 * @param action the action to be performed, if a value is present
 * @param emptyAction the empty-based action to be performed, if no value is
 *        present
 * @throws NullPointerException if a value is present and the given action
 *         is {@code null}, or no value is present and the given empty-based
 *         action is {@code null}.
 * @since 9
 */
func (o Optional) IfPresentOrElse(action consumer.Consumer, emptyAction EmptyConsumer) {
	if o.IsPresent() {
		action.Accept(o.value)
		return
	}
	emptyAction.Run()
}

/**
 * If a value is present, and the value matches the given predicate,
 * returns an {@code Optional} describing the value, otherwise returns an
 * empty {@code Optional}.
 *
 * @param predicate the predicate to apply to a value, if present
 * @return an {@code Optional} describing the value of this
 *         {@code Optional}, if a value is present and the value matches the
 *         given predicate, otherwise an empty {@code Optional}
 * @throws NullPointerException if the predicate is {@code null}
 */
func (o Optional) Filter(predicate predicate.Predicater) Optional {
	object.RequireNonNull(predicate)
	if !o.IsPresent() {
		return o
	}
	if predicate.Test(o.value) {
		return o
	}
	return Empty()
}

/**
 * If a value is present, returns an {@code Optional} describing (as if by
 * {@link #ofNullable}) the result of applying the given mapping function to
 * the value, otherwise returns an empty {@code Optional}.
 *
 * <p>If the mapping function returns a {@code null} result then this method
 * returns an empty {@code Optional}.
 *
 * @apiNote
 * This method supports post-processing on {@code Optional} values, without
 * the need to explicitly check for a return status.  For example, the
 * following code traverses a stream of URIs, selects one that has not
 * yet been processed, and creates a path from that URI, returning
 * an {@code Optional<Path>}:
 *
 * <pre>{@code
 *     Optional<Path> p =
 *         uris.stream().filter(uri -> !isProcessedYet(uri))
 *                       .findFirst()
 *                       .map(Paths::get);
 * }</pre>
 *
 * Here, {@code findFirst} returns an {@code Optional<URI>}, and then
 * {@code map} returns an {@code Optional<Path>} for the desired
 * URI if one exists.
 *
 * @param mapper the mapping function to apply to a value, if present
 * @param <U> The type of the value returned from the mapping function
 * @return an {@code Optional} describing the result of applying a mapping
 *         function to the value of this {@code Optional}, if a value is
 *         present, otherwise an empty {@code Optional}
 * @throws NullPointerException if the mapping function is {@code null}
 */
func (o Optional) Map(mapper func(interface{}) interface{}) Optional {
	object.RequireNonNil(mapper)
	if !o.IsPresent() {
		return Empty()
	}

	return OfNillable(mapper(o.value))
}

/**
 * If a value is present, returns the result of applying the given
 * {@code Optional}-bearing mapping function to the value, otherwise returns
 * an empty {@code Optional}.
 *
 * <p>This method is similar to {@link #map(Function)}, but the mapping
 * function is one whose result is already an {@code Optional}, and if
 * invoked, {@code flatMap} does not wrap it within an additional
 * {@code Optional}.
 *
 * @param <U> The type of value of the {@code Optional} returned by the
 *            mapping function
 * @param mapper the mapping function to apply to a value, if present
 * @return the result of applying an {@code Optional}-bearing mapping
 *         function to the value of this {@code Optional}, if a value is
 *         present, otherwise an empty {@code Optional}
 * @throws NullPointerException if the mapping function is {@code null} or
 *         returns a {@code null} result
 */
func (o *Optional) FlatMap(mapper func(interface{}) interface{}) Optional {
	object.RequireNonNil(mapper)
	if !o.IsPresent() {
		return Empty()
	}

	return OfNillable(mapper(o.value))
}
