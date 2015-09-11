package tmx_test

import (
	"os"

	. "github.com/manyminds/tmx"
	"github.com/manyminds/tmx/spec"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Loader", func() {
	Context("Load Flipped Tiles correctly", func() {
		var (
			testfile string
			testMap  *Map
		)

		BeforeEach(func() {
			testfile = "testfiles/simple_example.tmx"
			file, err := os.Open(testfile)
			Expect(err).ToNot(HaveOccurred())
			defer file.Close()

			testMap, err = NewMap(file)
			Expect(err).ToNot(HaveOccurred())
		})

		It("Should have at least one flipped tile", func() {
			Expect(len(testMap.Layers)).To(BeNumerically(">", 1))
			layer := testMap.Layers[1]
			horizontalFlippedCount := 0
			verticalFlippedCount := 0
			diagonallyFlippedCount := 0
			Expect(len(layer.Data.DataTiles)).To(BeNumerically(">", 1))
			for _, tile := range layer.Data.DataTiles {
				if tile.DiagonalFlip {
					diagonallyFlippedCount++
				}

				if tile.VerticalFlip {
					verticalFlippedCount++
				}

				if tile.HorizontalFlip {
					horizontalFlippedCount++
				}

				_, err := testMap.GetTilesetForGID(tile.GID)
				Expect(err).ToNot(HaveOccurred())
			}

			Expect(verticalFlippedCount).To(BeNumerically(">", 1))
			Expect(horizontalFlippedCount).To(BeNumerically(">", 1))
			Expect(diagonallyFlippedCount).To(BeNumerically(">", 1))
		})
	})

	Context("Load TMX Files", func() {
		It("Should load a simple valid file", func() {
			testfile := "testfiles/simple_example.tmx"
			file, err := os.Open(testfile)
			Expect(err).ToNot(HaveOccurred())
			defer file.Close()

			target, err := NewMap(file)
			Expect(err).ToNot(HaveOccurred())
			Expect(target.Tilesets).To(HaveLen(2))
			tileSet := target.Tilesets[0]

			Expect(string(target.BackgroundColor)).To(Equal("#248026"))
			Expect(target.Orientation).To(Equal("orthogonal"))
			Expect(target.Layers).To(HaveLen(3))

			Expect(target.Height).To(Equal(24))
			Expect(target.Width).To(Equal(24))
			Expect(target.RenderOrder).To(Equal("right-down"))
			Expect(tileSet.FirstGID).To(Equal(spec.GID(1)))
			Expect(tileSet.Name).To(Equal("chipset"))
			Expect(tileSet.TileHeight).To(Equal(32))
			Expect(tileSet.TileWidth).To(Equal(32))
			Expect(tileSet.Spacing).To(Equal(0))
			Expect(tileSet.Margin).To(Equal(0))

			//Tiles are only stores for tiles with specific inspecion
			//such as custom propertys
			//or animations
			Expect(tileSet.Tiles).To(HaveLen(2))

			image := tileSet.Image
			Expect(image.Source).To(Equal("chipset.png"))
			Expect(image.Width).To(Equal(320))
			Expect(image.Height).To(Equal(1600))
			floorLayer := target.Layers[0]
			belowLayer := target.Layers[1]
			aboveLayer := target.Layers[2]

			mapWidth := floorLayer.Width
			mapHeight := floorLayer.Height
			Expect(mapWidth).To(Equal(24))
			Expect(mapHeight).To(Equal(24))

			Expect(floorLayer.Width).To(Equal(belowLayer.Width))
			Expect(aboveLayer.Width).To(Equal(belowLayer.Width))

			Expect(floorLayer.Data.DataTiles).To(HaveLen(mapWidth * mapHeight))
			for _, tile := range floorLayer.Data.DataTiles {
				Expect(target.GetTilesetForGID(tile.GID)).To(Equal(&target.Tilesets[0]))
			}

			Expect(floorLayer.IsVisible()).To(Equal(true))
			Expect(belowLayer.IsVisible()).To(Equal(true))
			Expect(aboveLayer.IsVisible()).To(Equal(true))
		})
	})
})
