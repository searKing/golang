// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package spliterator

type Characteristic int

const (
	CharacteristicTODO Characteristic = 0x00000000
	/**
	 * Characteristic value signifying that an encounter order is defined for
	 * elements. If so, this Spliterator guarantees that method
	 * {@link #trySplit} splits a strict prefix of elements, that method
	 * {@link #tryAdvance} steps by one element in prefix order, and that
	 * {@link #forEachRemaining} performs actions in encounter order.
	 *
	 * <p>A {@link Collection} has an encounter order if the corresponding
	 * {@link Collection#iterator} documents an order. If so, the encounter
	 * order is the same as the documented order. Otherwise, a collection does
	 * not have an encounter order.
	 *
	 * @apiNote Encounter order is guaranteed to be ascending index order for
	 * any {@link List}. But no order is guaranteed for hash-based collections
	 * such as {@link HashSet}. Clients of a Spliterator that reports
	 * {@code ORDERED} are expected to preserve ordering constraints in
	 * non-commutative parallel computations.
	 */
	CharacteristicOrdered Characteristic = 0x00000010

	/**
	 * Characteristic value signifying that, for each pair of
	 * encountered elements {@code x, y}, {@code !x.equals(y)}. This
	 * applies for example, to a Spliterator based on a {@link Set}.
	 */
	CharacteristicDistinct Characteristic = 0x00000001

	/**
	 * Characteristic value signifying that encounter order follows a defined
	 * sort order. If so, method {@link #getComparator()} returns the associated
	 * Comparator, or {@code null} if all elements are {@link Comparable} and
	 * are sorted by their natural ordering.
	 *
	 * <p>A Spliterator that reports {@code SORTED} must also report
	 * {@code ORDERED}.
	 *
	 * @apiNote The spliterators for {@code Collection} classes in the JDK that
	 * implement {@link NavigableSet} or {@link SortedSet} report {@code SORTED}.
	 */
	CharacteristicSorted Characteristic = 0x00000004

	/**
	 * Characteristic value signifying that the value returned from
	 * {@code estimateSize()} prior to traversal or splitting represents a
	 * finite size that, in the absence of structural source modification,
	 * represents an exact count of the number of elements that would be
	 * encountered by a complete traversal.
	 *
	 * @apiNote Most Spliterators for Collections, that cover all elements of a
	 * {@code Collection} report this characteristic. Sub-spliterators, such as
	 * those for {@link HashSet}, that cover a sub-set of elements and
	 * approximate their reported size do not.
	 */
	CharacteristicSized Characteristic = 0x00000040

	/**
	 * Characteristic value signifying that the source guarantees that
	 * encountered elements will not be {@code null}. (This applies,
	 * for example, to most concurrent collections, queues, and maps.)
	 */
	CharacteristicNonnulL Characteristic = 0x00000100

	/**
	 * Characteristic value signifying that the element source cannot be
	 * structurally modified; that is, elements cannot be added, replaced, or
	 * removed, so such changes cannot occur during traversal. A Spliterator
	 * that does not report {@code IMMUTABLE} or {@code CONCURRENT} is expected
	 * to have a documented policy (for example throwing
	 * {@link ConcurrentModificationException}) concerning structural
	 * interference detected during traversal.
	 */
	CharacteristicImmutable Characteristic = 0x00000400

	/**
	 * Characteristic value signifying that the element source may be safely
	 * concurrently modified (allowing additions, replacements, and/or removals)
	 * by multiple threads without external synchronization. If so, the
	 * Spliterator is expected to have a documented policy concerning the impact
	 * of modifications during traversal.
	 *
	 * <p>A top-level Spliterator should not report both {@code CONCURRENT} and
	 * {@code SpliteratorSIZED}, since the finite size, if known, may change if the source
	 * is concurrently modified during traversal. Such a Spliterator is
	 * inconsistent and no guarantees can be made about any computation using
	 * that Spliterator. Sub-spliterators may report {@code SpliteratorSIZED} if the
	 * sub-split size is known and additions or removals to the source are not
	 * reflected when traversing.
	 *
	 * <p>A top-level Spliterator should not report both {@code CONCURRENT} and
	 * {@code IMMUTABLE}, since they are mutually exclusive. Such a Spliterator
	 * is inconsistent and no guarantees can be made about any computation using
	 * that Spliterator. Sub-spliterators may report {@code IMMUTABLE} if
	 * additions or removals to the source are not reflected when traversing.
	 *
	 * @apiNote Most concurrent collections maintain a consistency policy
	 * guaranteeing accuracy with respect to elements present at the point of
	 * Spliterator construction, but possibly not reflecting subsequent
	 * additions or removals.
	 */
	CharacteristicConcurrent Characteristic = 0x00001000

	/**
	 * Characteristic value signifying that all Spliterators resulting from
	 * {@code trySplit()} will be both {@link #SpliteratorSIZED} and {@link #SUBSIZED}.
	 * (This means that all child Spliterators, whether direct or indirect, will
	 * be {@code SpliteratorSIZED}.)
	 *
	 * <p>A Spliterator that does not report {@code SpliteratorSIZED} as required by
	 * {@code SUBSIZED} is inconsistent and no guarantees can be made about any
	 * computation using that Spliterator.
	 *
	 * @apiNote Some spliterators, such as the top-level spliterator for an
	 * approximately balanced binary tree, will report {@code SpliteratorSIZED} but not
	 * {@code SUBSIZED}, since it is common to know the size of the entire tree
	 * but not the exact sizes of subtrees.
	 */
	CharacteristicSubsized Characteristic = 0x00004000
)
